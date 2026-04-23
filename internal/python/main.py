import ee
import requests
from geopy.geocoders import Nominatim
from PIL import Image
from io import BytesIO
import sys

ee.Initialize(project='sentinel-data-492807')

class Handler():
    @staticmethod
    def GetCoordinates(city_name):
        geolocator = Nominatim(user_agent="geo_app")
        location = geolocator.geocode(city_name)

        if not location:
            raise ValueError("City not found")

        return location.latitude, location.longitude

    @staticmethod
    def CreateBBox(lat, lon, delta=0.05):
        return ee.Geometry.Rectangle([
            lon - delta,
            lat - delta,
            lon + delta,
            lat + delta
        ])

class Service():
    def addCloudScore(self, image):
        scl = image.select('SCL')

        cloud_mask = (
            scl.eq(3)   # тени
            .Or(scl.eq(8))  # облака
            .Or(scl.eq(9))
            .Or(scl.eq(10))
            .Or(scl.eq(11))  # снег (!!)
        )

        cloud_fraction = cloud_mask.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=60,
            maxPixels=1e9
        ).get('SCL')

        cloud_fraction = ee.Number(
            ee.Algorithms.If(cloud_fraction, cloud_fraction, 1)
        )

        return image.set('CLOUD_SCORE_BBOX', cloud_fraction)
    
    def addCoverage(self, image):
        valid_mask = image.select('B4').mask()

        coverage = valid_mask.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=20,
            maxPixels=1e9
        ).get('B4')

        return image.set('COVERAGE', coverage)

    def addWhitenessScore(self, image):
        r = image.select('B4')
        g = image.select('B3')
        b = image.select('B2')

        mean = r.add(g).add(b).divide(3)

        whiteness = (
            r.subtract(mean).abs()
            .add(g.subtract(mean).abs())
            .add(b.subtract(mean).abs())
        )

        white_pixels = mean.gt(2000).And(whiteness.lt(200))

        white_fraction = white_pixels.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=60,
            maxPixels=1e9
        ).get('B4')

        white_fraction = ee.Number(
            ee.Algorithms.If(white_fraction, white_fraction, 1)
        )

        return image.set('WHITE_SCORE', white_fraction)

    def GetImage(self):
        if len(sys.argv) < 2:
            raise ValueError("City argument missing")
        
        self.city = sys.argv[1].split("=")[-1]
        self.date_from = sys.argv[2].split("=")[-1]
        self.date_to = sys.argv[3].split("=")[-1]

        self.lat, self.lon = Handler.GetCoordinates(self.city)
        self.bbox = Handler.CreateBBox(self.lat, self.lon)

        collection = (
            ee.ImageCollection("COPERNICUS/S2_SR_HARMONIZED")
            .filterBounds(self.bbox)
            .filterDate(self.date_from, self.date_to)
            .filter(ee.Filter.lt("CLOUDY_PIXEL_PERCENTAGE", 50))
            .limit(20)
            .map(self.addCloudScore)
            .map(self.addWhitenessScore)
            .map(self.addCoverage)
            .filter(ee.Filter.gt("COVERAGE", 0.9))
            .sort("WHITE_SCORE")
            .sort("CLOUD_SCORE_BBOX")
        )

        size = collection.size().getInfo()
        if size == 0:
            raise Exception("No images found for given filters")

        image = collection.limit(5).median()

        url = image.getThumbURL({
            "region": self.bbox,
            "dimensions": 512,
            "format": "png",
            "bands": ["B4", "B3", "B2"],
            "min": 0,
            "max": 3000
        })

        response = requests.get(url)
        if response.status_code != 200:
            raise Exception("Failed to fetch image")
        
        if response.status_code != 200:
            raise Exception(f"HTTP error: {response.status_code}")

        content_type = response.headers.get("Content-Type", "")

        if not content_type.startswith("image"):
            raise Exception(f"Not an image! Got: {content_type}\n{response.text[:200]}")

        sys.stdout.buffer.write(response.content)


if __name__ == "__main__":
    service = Service()
    service.GetImage()
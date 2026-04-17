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

        cloud_mask = scl.eq(8).Or(scl.eq(9)).Or(scl.eq(10))

        cloud_fraction = cloud_mask.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=20,
            maxPixels=1e9
        ).get('SCL')

        return image.set('CLOUD_SCORE_BBOX', cloud_fraction)

    def GetImage(self):
        if len(sys.argv) < 2:
            raise ValueError("City argument missing")
        
        self.city = sys.argv[1].split("=")[-1]

        self.lat, self.lon = Handler.GetCoordinates(self.city)
        self.bbox = Handler.CreateBBox(self.lat, self.lon)

        collection = (
            ee.ImageCollection("COPERNICUS/S2_SR_HARMONIZED")
            .filterBounds(self.bbox)
            .filterDate("2022-06-01", "2024-06-30")
            .map(self.addCloudScore)
            .sort("CLOUD_SCORE_BBOX")
            .filter(ee.Filter.lt("CLOUDY_PIXEL_PERCENTAGE", 80))
        )

        image = ee.Image(collection.first())

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
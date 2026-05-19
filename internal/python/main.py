from email.mime import image

import ee
import requests
from geopy.geocoders import Nominatim
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
            scl.eq(3)   # cloud shadow
            .Or(scl.eq(8))   # cloud medium probability
            .Or(scl.eq(9))   # cloud high probability
            .Or(scl.eq(10))  # cirrus
            .Or(scl.eq(11))  # snow
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

        return image.set(
            'CLOUD_SCORE_BBOX',
            cloud_fraction
        )

    def addCoverage(self, image):
        valid_mask = image.select(self.bands[0]).mask()

        coverage = valid_mask.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=20,
            maxPixels=1e9
        ).get(self.bands[0])

        coverage = ee.Number(
            ee.Algorithms.If(coverage, coverage, 0)
        )

        return image.set(
            'COVERAGE',
            coverage
        )

    def addWhitenessScore(self, image):

        # whiteness meaningful only for RGB-like combinations
        if set(self.bands) != set(['B4', 'B3', 'B2']):
            return image.set('WHITE_SCORE', 0)

        r = image.select('B4')
        g = image.select('B3')
        b = image.select('B2')

        mean = r.add(g).add(b).divide(3)

        whiteness = (
            r.subtract(mean).abs()
            .add(g.subtract(mean).abs())
            .add(b.subtract(mean).abs())
        )

        white_pixels = (
            mean.gt(1800)
            .And(whiteness.lt(250))
        )

        white_fraction = white_pixels.reduceRegion(
            reducer=ee.Reducer.mean(),
            geometry=self.bbox,
            scale=60,
            maxPixels=1e9
        ).get('B4')

        white_fraction = ee.Number(
            ee.Algorithms.If(
                white_fraction,
                white_fraction,
                1
            )
        )

        return image.set(
            'WHITE_SCORE',
            white_fraction
        )

    def addContrastScore(self, image):

        stats = image.select(
            self.bands
        ).reduceRegion(
            reducer=ee.Reducer.stdDev(),
            geometry=self.bbox,
            scale=60,
            maxPixels=1e9
        )

        values = []

        for band in self.bands:
            value = ee.Number(
                ee.Algorithms.If(
                    stats.get(band),
                    stats.get(band),
                    0
                )
            )

            values.append(value)

        contrast = ee.List(values).reduce(
            ee.Reducer.mean()
        )

        return image.set(
            'CONTRAST_SCORE',
            contrast
        )

    def addFinalScore(self, image):

        cloud = ee.Number(
            image.get('CLOUD_SCORE_BBOX')
        )

        white = ee.Number(
            image.get('WHITE_SCORE')
        )

        contrast = ee.Number(
            image.get('CONTRAST_SCORE')
        )

        final_score = (
            cloud.multiply(0.65)
            .add(white.multiply(0.25))
            .subtract(
                contrast.multiply(0.00003)
            )
        )

        return image.set(
            'FINAL_SCORE',
            final_score
        )

    def buildCollection(
        self,
        cloud_limit,
        white_limit
    ):

        return (
            ee.ImageCollection(
                "COPERNICUS/S2_SR_HARMONIZED"
            )

            .filterBounds(self.bbox)

            .filterDate(
                self.date_from,
                self.date_to
            )

            .filter(
                ee.Filter.lt(
                    "CLOUDY_PIXEL_PERCENTAGE",
                    50
                )
            )

            .limit(50)

            .map(self.addCloudScore)

            .map(self.addWhitenessScore)

            .map(self.addCoverage)

            .map(self.addContrastScore)

            .map(self.addFinalScore)

            .filter(
                ee.Filter.gt(
                    "COVERAGE",
                    0.9
                )
            )

            .filter(
                ee.Filter.lt(
                    "CLOUD_SCORE_BBOX",
                    cloud_limit
                )
            )

            .filter(
                ee.Filter.lt(
                    "WHITE_SCORE",
                    white_limit
                )
            )

            .sort("FINAL_SCORE")
        )

    def GetImage(self):

        if len(sys.argv) < 6:
            raise ValueError(
                "Arguments missing: "
                "city=..., "
                "date_from=..., "
                "date_to=..., "
                "bands=B4,B3,B2, "
                "dimensions=512"
            )

        self.city = (
            sys.argv[1]
            .split("=")[-1]
        )

        self.date_from = (
            sys.argv[2]
            .split("=")[-1]
        )

        self.date_to = (
            sys.argv[3]
            .split("=")[-1]
        )

        self.bands = (
            sys.argv[4]
            .split("=")[-1]
            .split(",")
        )

        self.dimensions = int(
            sys.argv[5]
            .split("=")[-1]
        )

        self.scale = int(
            sys.argv[6]
            .split("=")[-1]
        )

        self.output_format = (
            sys.argv[7]
            .split("=")[-1]
            .lower()
        )

        if self.output_format not in [
            "png",
            "geotiff"
        ]:
            raise ValueError(
                "format must be png or geotiff"
            )

        if len(self.bands) != 3:
            raise ValueError(
                "Exactly 3 bands required"
            )

        self.lat, self.lon = (
            Handler.GetCoordinates(
                self.city
            )
        )

        self.bbox = Handler.CreateBBox(
            self.lat,
            self.lon
        )

        fallback_configs = [
            (0.05, 0.05),
            (0.08, 0.08),
            (0.12, 0.12),
            (0.18, 0.18),
            (0.25, 0.25),
        ]

        collection = None

        for cloud_limit, white_limit in fallback_configs:

            candidate = self.buildCollection(
                cloud_limit,
                white_limit
            )

            size = candidate.size().getInfo()

            if size > 0:
                collection = candidate
                break

        if collection is None:
            raise Exception(
                "No suitable images found"
            )

        if self.output_format == "png":
            image = (
                collection
                .limit(5)
                .median()
            )

            url = image.getThumbURL({
                "region": self.bbox,
                "dimensions": self.dimensions,
                "format": "png",
                "bands": self.bands,
                "min": 0,
                "max": 3000
            })

        else:
            image = (
                collection
                .first()
            )

            url = image.getDownloadURL({
                "region": self.bbox,
                "scale": self.scale,
                "bands": self.bands,
                "format": "GEO_TIFF",
                "crs": "EPSG:4326"
            })

        response = requests.get(url)

        if response.status_code != 200:
            raise Exception(
                f"HTTP error: {response.status_code}"
            )

        sys.stdout.buffer.write(response.content)


if __name__ == "__main__":
    service = Service()
    service.GetImage()
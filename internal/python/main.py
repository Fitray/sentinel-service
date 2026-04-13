import ee
import requests
from geopy.geocoders import Nominatim
from PIL import Image
from io import BytesIO

# Инициализация GEE
ee.Authenticate()
ee.Initialize(project='sentinel-data-492807')


class Handler():
    @staticmethod
    def GetCoordinates(city_name):
        geolocator = Nominatim(user_agent="geo_app")
        location = geolocator.geocode(city_name)

        if not location:
            raise ValueError("Населённый пункт не найден")

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
    def __init__(self, bbox):
        self.bbox = bbox

    def DoRequest(self):
        # Коллекция Sentinel-2
        collection = (
            ee.ImageCollection("COPERNICUS/S2_SR_HARMONIZED")
            .filterBounds(self.bbox)
            .filterDate("2024-06-01", "2024-06-30")
            .filter(ee.Filter.lt("CLOUDY_PIXEL_PERCENTAGE", 20))
        )

        image = collection.sort("CLOUDY_PIXEL_PERCENTAGE").first()

        # RGB визуализация (аналог evalscript)
        vis_params = {
            "bands": ["B4", "B3", "B2"],
            "min": 0,
            "max": 3000
        }

        url = image.getThumbURL({
            "region": self.bbox,
            "dimensions": 512,
            "format": "png",
            **vis_params
        })

        # Скачиваем изображение
        response = requests.get(url)
        if response.status_code != 200:
            raise ValueError("Ошибка при загрузке изображения")

        img = Image.open(BytesIO(response.content))
        img.save("result.png")

        print("Изображение сохранено как result.png")


# === Запуск ===

city = input("Укажите территорию: ")

lat, lon = Handler.GetCoordinates(city)
bbox = Handler.CreateBBox(lat, lon)

service = Service(bbox=bbox)
service.DoRequest()
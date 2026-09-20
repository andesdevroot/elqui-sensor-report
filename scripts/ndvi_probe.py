#!/usr/bin/env python3
"""
ndvi_probe.py — Consulta Earth Search y calcula NDVI de una parcela.

Uso:
    python3 scripts/ndvi_probe.py --lat -29.90453 --lon -71.24894 --days 90 --cloud 20

Sin tokens, sin OAuth. Earth Search es un STAC público anónimo.
"""
import argparse
from datetime import datetime, timedelta, timezone

import numpy as np
import rasterio
from rasterio.warp import transform as warp_transform
from rasterio.windows import from_bounds
from pystac_client import Client

EARTH_SEARCH = "https://earth-search.aws.element84.com/v1"
COLLECTION = "sentinel-2-l2a"


def find_items(lat: float, lon: float, days: int, max_cloud: float):
    """Busca items Sentinel-2 L2A recientes cerca de (lat, lon)."""
    end = datetime.now(timezone.utc)
    start = end - timedelta(days=days)

    client = Client.open(EARTH_SEARCH)
    search = client.search(
        collections=[COLLECTION],
        intersects={"type": "Point", "coordinates": [lon, lat]},
        datetime=f"{start.isoformat()}/{end.isoformat()}",
        query={"eo:cloud_cover": {"lt": max_cloud}},
        max_items=20,
    )
    return sorted(search.items(), key=lambda i: i.datetime, reverse=True)


def read_band_window(item, band_key: str, lat: float, lon: float, half_meters: float = 100.0):
    """
    Lee una ventana cuadrada de ~half_meters de lado alrededor de (lat, lon).

    Devuelve (array, transform, crs). Si la banda tiene resolución distinta
    (10m vs 20m), devuelve arrays de distinto tamaño — se maneja en el caller.
    """
    asset = item.assets[band_key]
    with rasterio.open(asset.href) as src:
        xs, ys = warp_transform("EPSG:4326", src.crs, [lon], [lat])
        x, y = xs[0], ys[0]

        window = from_bounds(
            x - half_meters, y - half_meters,
            x + half_meters, y + half_meters,
            transform=src.transform,
        )
        data = src.read(1, window=window)
        return data, src.transform, src.crs


def compute_ndvi(red: np.ndarray, nir: np.ndarray) -> np.ndarray:
    """NDVI puro. Ignora máscara de nubes por ahora (Fase 1.2)."""
    red = red.astype(np.float32)
    nir = nir.astype(np.float32)
    denom = nir + red
    return np.where(denom > 0, (nir - red) / denom, np.nan)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--lat", type=float, required=True)
    parser.add_argument("--lon", type=float, required=True)
    parser.add_argument("--days", type=int, default=90)
    parser.add_argument("--cloud", type=float, default=20)
    parser.add_argument("--limit", type=int, default=3, help="Cuántos items procesar")
    args = parser.parse_args()

    print(f"Buscando imágenes para ({args.lat}, {args.lon}), últimos {args.days} días, cloud<{args.cloud}%")
    items = find_items(args.lat, args.lon, args.days, args.cloud)
    print(f"Encontradas {len(items)} imágenes\n")
    if not items:
        print("Sin imágenes. Amplía --days o --cloud.")
        return

    for item in items[: args.limit]:
        cloud = item.properties.get("eo:cloud_cover", -1)
        print(f"--- {item.id} | {item.datetime.date()} | cloud={cloud:.2f}% ---")
        try:
            red, _, _ = read_band_window(item, "red", args.lat, args.lon)
            nir, _, _ = read_band_window(item, "nir", args.lat, args.lon)

            if red.size == 0 or nir.size == 0:
                print("  Ventana vacía (punto fuera del raster)")
                continue

            if red.shape != nir.shape:
                print(f"  Shapes incompatibles: red={red.shape}, nir={nir.shape}")
                continue

            ndvi = compute_ndvi(red, nir)
            valid = ndvi[np.isfinite(ndvi)]
            if valid.size == 0:
                print("  Sin píxeles válidos")
                continue

            print(f"  Píxeles: {valid.size}")
            print(f"  NDVI: media={valid.mean():.3f} min={valid.min():.3f} max={valid.max():.3f}")
        except Exception as e:
            print(f"  Error: {type(e).__name__}: {e}")
        print()


if __name__ == "__main__":
    main()

    
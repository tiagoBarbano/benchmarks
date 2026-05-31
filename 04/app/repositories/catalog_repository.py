from pathlib import Path

import yaml


class CatalogRepository:

    def __init__(self):

        self._file = Path(
            "config/apis.yaml"
        )

    def load(self):

        with open(self._file) as f:

            return yaml.safe_load(f)
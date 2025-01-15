import json
import numpy as np
import pandas as pd

import settings

DATA_LOCATION =  '/Users/brendan/Documents/projects/hotspot/go/data/'

COLUMN_NAMES = {
    "RunSets.Parameters.R0": "R0",
    "HotspotFraction": "Hotspot fraction",
    "RiskVariance": "Risk tolerance variance",
    # make sure RiskMean is made categorical first
    "RiskMean": "Risk tolerance mean",
    "MaxIDiff": "Peak size difference",
    "FinalRDiff": "Extent difference",
    "OutbreakProbability": "Outbreak probability",
}

# Cache of data. Only really does anything if we're using jupyter to keep
# the script alive but it could be nice.
DATA_CACHE = {}


def load_data(filename, drop_control=True, reload=False):
    if filename in DATA_CACHE and not reload:
        return DATA_CACHE[filename]

    with open(DATA_LOCATION + filename) as file:
            json_file = json.load(file, parse_float=lambda f: round(float(f), 2))
            
    data = process(pd.json_normalize(
        json_file,
        record_path=["RunSets", "Runs"],
        meta=[
            "RiskMean",
            "RiskVariance",
            "HotspotFraction",
            ["RunSets", "Parameters", "R0"],
            ["RunSets", "Parameters", "RunType"]
        ],
    ), drop_control=drop_control)

    DATA_CACHE[filename] = data

    return data


# Process data
def process(data, drop_control):
    
    data["OutbreakProbability"] = 1.0 - 1.0*(data["FinalR"] <= settings.EXTINCTION_CUTOFF)
    data["RiskMean"] = pd.Categorical(data["RiskMean"])
    
    ## Uses np.tile to replicate the control series
    def tile(column):
        num_ppf = len(data["HotspotFraction"].unique())
        return pd.Series(np.tile(data[data["HotspotFraction"] == 0][column], num_ppf))
    
    data["MaxIControl"] = tile("MaxI")
    data["FinalRControl"] =  tile("FinalR")
    
    data["FinalRDiff"] = data["FinalR"] - data["FinalRControl"]
    data["MaxIDiff"] = data["MaxI"] - data["MaxIControl"]

    if drop_control:
        data = data[data["HotspotFraction"] != 0]
    
    data = data.rename(columns=COLUMN_NAMES)
    return data


def set_width(plot):
    w, h = plot.figure.get_size_inches()
    plot.figure.set_size_inches(settings.FULL_WIDTH, h)
import json
import numpy as np
import pandas as pd

DATA_LOCATION =  '/Users/brendan/Documents/projects/hotspot/go/data/'

COLUMN_NAMES = {
    "RunSets.Parameters.R0": "$R_0$",
    "HotspotFraction": "Hotspot fraction",
    "RiskVariance": "Risk tolerance variance",
    # make sure RiskMean is made categorical first
    "RiskMean": "Risk tolerance mean",
    "MaxIDiff": "Peak size difference",
    "FinalRDiff": "Extent difference",
    "OutbreakProbability": "Outbreak probability",
}


def load_data(file_name):

    with open(DATA_LOCATION + file_name) as file:
            json_file = json.load(file, parse_float=lambda f: round(float(f), 2))
            
    data = pd.json_normalize(
        json_file,
        record_path=["RunSets", "Runs"],
        meta=[
            "RiskMean",
            "RiskVariance",
            "HotspotFraction",
            ["RunSets", "Parameters", "R0"],
            ["RunSets", "Parameters", "RunType"]
        ],
    )
    return data

# data = load_data("simulation,D=1,T=1000.json")


# Process data
def process_data(data):
    
    data["OutbreakProbability"] = 1.0 - 1.0*(data["FinalR"] <= EXTINCTION_CUTOFF)
    data["RiskMean"] = pd.Categorical(data["RiskMean"])
    
    ## Uses np.tile to replicate the control series
    def tile(column):
        num_ppf = len(data["HotspotFraction"].unique())
        return pd.Series(np.tile(data[data["HotspotFraction"] == 0][column], num_ppf))
    
    data["MaxIControl"] = tile("MaxI")
    data["FinalRControl"] =  tile("FinalR")
    
    data["FinalRDiff"] = data["FinalR"] - data["FinalRControl"]
    data["MaxIDiff"] = data["MaxI"] - data["MaxIControl"]
    data = data[data["HotspotFraction"] != 0]
    
    data = data.rename(columns=COLUMN_NAMES)
    return data

#!/usr/bin/python3

import json
import sys

if __name__ == "__main__":
    #json_file = "config_db.json.l3vlan"
    json_file = sys.argv[1]

    f = open(json_file)

    try:
        json.load(f)
    except ValueError as err:
        print("Error. Check {}".format(json_file))
    else:
        print("ok")

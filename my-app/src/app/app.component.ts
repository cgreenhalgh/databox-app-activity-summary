import { Component } from '@angular/core';
import { OnInit } from '@angular/core';

import { DataService } from './data.service';
import { DataType } from './types';

import * as vega from 'vega';
import * as vegaLite from 'vega-lite';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  loading = true;
  dataTypes: DataType[];
  constructor(private dataService: DataService) {}
  ngOnInit(): void {
	  this.dataService.getDataTypes().then(dts => { this.dataTypes = dts; this.loading = false });
  
    // test example
    let vlSpec = getExVegaLiteSpec();
    let theSpec = vegaLite.compile(vlSpec).spec;
   
    const container = document.querySelector("#view");
    const runtime = vega.parse(theSpec, {});
    const view = new vega.View(runtime, { logLevel: vega.Warn });
    const initializedView = view.initialize(container);
    view.run();
  }
}
function getExVegaLiteSpec(): any {
  return {
    "$schema": "https://vega.github.io/schema/vega-lite/v2.json",
    "description": "A simple bar chart with embedded data.",
    "width": 360,
    "data": { 
      "url": "../api/get/strava_activity/"
    /*"values": [
        {"a": "A","b": 28 },
        {"a": "B","b": 55},
        {"a": "C","b": 43},
        {"a": "D","b": 91},
        {"a": "E","b": 81},
        {"a": "F","b": 53},
        {"a": "G","b": 19},
        {"a": "H","b": 87},
        {"a": "I","b": 52}
      ]*/
    },
    "mark": "bar",
    "transform": [
      {"calculate": "datum.data.distance*0.001", "as": "distance"},
      {"calculate": "floor(datum.timestamp / (24*60*60*1000))*(24*60*60*1000)", "as": "date"},
    ],
    "encoding": {
      "x": {
        "field": "date",
        "type": "temporal"
        /*,"timeUnit": "day"*/
      },
      "y": {
        "field": "distance",
        "type": "quantitative",
        "aggregate": "sum",
        "axis":{
          "title":"Total distance (km)"
        }
      }
    }
  }
}

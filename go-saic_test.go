package main

import "testing"

type calcSizeParms struct {
	imageWidth        int
	imageHeight       int
	tileWidth         int
	tileHeight        int
	expectedHorzTiles int
	expectedVertTiles int
}

func TestCalcTileSize(t *testing.T) {

	testParmsList := map[string]calcSizeParms{
		"square image": calcSizeParms{imageWidth: 128,
			imageHeight:       128,
			tileWidth:         64,
			tileHeight:        64,
			expectedHorzTiles: 2,
			expectedVertTiles: 2},
		"tall image": calcSizeParms{imageWidth: 64,
			imageHeight:       128,
			tileWidth:         64,
			tileHeight:        64,
			expectedHorzTiles: 1,
			expectedVertTiles: 2},
		"wide image": calcSizeParms{imageWidth: 128,
			imageHeight:       64,
			tileWidth:         64,
			tileHeight:        64,
			expectedHorzTiles: 2,
			expectedVertTiles: 1},
		"tall image with overhang": calcSizeParms{imageWidth: 64,
			imageHeight:       130,
			tileWidth:         64,
			tileHeight:        64,
			expectedHorzTiles: 1,
			expectedVertTiles: 3},
		"wide image with overhang": calcSizeParms{imageWidth: 130,
			imageHeight:       64,
			tileWidth:         64,
			tileHeight:        64,
			expectedHorzTiles: 3,
			expectedVertTiles: 1},
	}

	for testName, testParms := range testParmsList {
		t.Logf("Running test:%s", testName)
		horzTiles, vertTiles := calcMosaicTiles(testParms.imageWidth, testParms.imageHeight, testParms.tileWidth, testParms.tileHeight)
		if horzTiles != testParms.expectedHorzTiles || vertTiles != testParms.expectedVertTiles {
			t.Errorf("Expected %v,%v Received %v,%v", testParms.expectedHorzTiles, testParms.expectedVertTiles, horzTiles, vertTiles)
		}
	}

}

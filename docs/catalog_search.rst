Catalog Searches
================

Catalog Search Overview
-----------------------

The catalog can be searched various ways.  Input to gbdxcli can be a line separated file, a list of space seperated parameters (startDate endDate type searchAreaWkt filter) or parameters entered one per line.  The JSON output can be formatted using unix command line tools.

.. code-block:: bash

   $ gbdxcli catalog < aoi_search.txt

The results will be a list of items, each of which contains quite a lot of metadata and looks something like this:

.. code-block:: json

	{
		"identifier": "1030010064562600",
		"type":["GBDXCatalogRecord","Acquisition","DigitalGlobeAcquisition","WV02"],
		"properties": {
		  "bearing":0,
		  "browseURL":"https://geobigdata.io/thumbnails/v1/browse/1030010064562600",
		  "catalogID":"1030010064562600",
		  "cloudCover":3,
		  "footprintWkt":"MULTIPOLYGON(((172.6160617 -78.2577799,...)))",
		  "imageBands":"PAN_MS1_MS2",
		  ...
		  "timestamp":"2017-01-17T14:35:08.000Z",
		  "vendor":"DigitalGlobe"
		}
	}

You could get a list of catalog IDs for this search by doing the following.  The JSON output is formatted using the jq utility which is a light-weight, multi-purpose utility available in Linux distros and Mac OSX.

.. code-block:: bash
	
	$ gbdxcli catalog < aoi_search.txt | jq ".[] | .identifier"

        "1020010052CCDB00"
        "101001000281E600"
        "103001005BA65A00"

Search by AOI
-----------------------
Search the catalog by AOI, as defined by a WKT polygon.  All imagery that intersects the polygon will be returned.

.. code-block::  bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\)


Search by Dates
-----------------------
The catalog can also be searched by date.  Note that if no search-polygon is supplied, the catalog only supports 
date searches of of one week intervals at a time.


.. code-block:: bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition

Search with Filters
-----------------------
You can add filters for any properties in the catalog items you are searching for.  For example, here's how you return images with cloudCover less than 50.:

.. code-block:: bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\) cloudCover \< 50


Here's a more complicated set of filters that can be applied:

.. code-block:: bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\) offNadirAngle \< 25,cloudCover \< 30

Search by Types
-----------------------
You can search by type as well.  The usual type for Digital Globe Imagery is "DigitalGlobeAcquisition".  
To search only Landsat imagery for example:

.. code-block:: bash

        $ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z LandsatAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\)

To search only IDAHO images example:

.. code-block:: bash

	$ gbdxcli catalog 2017-01-21T00:00:00.000Z 2017-01-21T23:59:59.000Z IDAHOIMAGE

Get Metadata Info about a given Catalog ID
------------------------------------------
If you have multiple catalog IDs and simply want to get records out of the catalog:
To get records by DigitalGlobe catalog ID:

.. code-block:: bash

         $ gbdxcli catalog get 103001005DB90000 103001005BA65A00 103001005BA6CC00

To get IDAHOImage records:

.. code-block:: bash

	$ gbdxcli catalog get 345b6306-60df-4180-a24c-6497be24b309


        {
		"identifier":"345b6306-60df-4180-a24c-6497be24b309",
		"type":["GBDXCatalogRecord","IDAHOImage"],
		"properties": {
		  "DGCatalogId":"7fdc5e35-abe9-4fa1-9362-fe335491a44d",
		  "bucketName":"idaho-images",
		  "colorInterpretation":"PAN",
		  "dataType":"UNSIGNED_SHORT",
		  "epsgCode":"4326",
		  "footprintWkt":"MULTIPOLYGON(((-61.82190258 -32.00228972, ...)))",
		  "groundSampleDistanceMeters":0.533,
		  "imageHeight":28068,
		  "imageId":"7fdc5e35-abe9-4fa1-9362-fe335491a44d",
		  "imageWidth":35180,
		  "nativeTileFileFormat":"PNG",
		  "numBands":1,
		  "numXTiles":35,
		  "numYTiles":28,
		  "platformName":"WORLDVIEW02",
		  "pniirs":4.9,
		  "profileName":"dg_1b",
		  "satElevation":66.0,
		  "tileBucketName":"idaho-images",
		  "tilePartition":"0000",
		  "tileXOffset":0,
		  "tileXSize":1024,
		  "tileYOffset":0,
		  "tileYSize":1024,
		  "timestamp":"2017-01-23T23:50:39.000Z",
		  "vendor":"DigitalGlobe",
		  "vendorDatasetIdentifier":"LV1B:056117417010_01_P010:1030010011348C00:A010010252754400",
		  "vendorDatasetIdentifier1":"LV1B",
	  	  "vendorDatasetIdentifier2":"056117417010_01_P010",
		  "vendorDatasetIdentifier3":"1030010011348C00",
		  "vendorDatasetIdentifier4":"A010010252754400",
		  "vendorName":"DigitalGlobe",
		  "version":"1.0"
		}
	}

To get LandSat records:

.. code-block:: bash

	$ gbdxcli catalog get LC80170402016246LGN00

Search File for a given search
------------------------------

You can redirect and reuse searches in a test file with line seperated values.

.. code-block:: bash

         $ cat aoi_search.txt

         2016-09-01T00:00:00.000Z
         2016-09-07T00:00:00.000Z
         DigitalGlobeAcquisition
         POLYGON((-113.88427734375 40.36642741921034,-110.28076171875 40.36642741921034,-110.28076171875 37.565262680889965,-113.88427734375 37.565262680889965,-113.88427734375 40.36642741921034))

	$ gbdxcli catalog < aoi_search.txt

Search for the most recent
--------------------------
You can search for the most recent and a limit number of records by adding "limit=" to the filter line.

.. code-block:: bash

	$ cat aoi_1B_wv03filter.txt

	2017-01-07T00:00:00.000Z
	2017-01-17T23:59:59.000Z
	1BProduct
	platformName = WORLDVIEW03,limit=recent

	$ gbdxcli catalog < aoi_1B_wv03filter.txt | wc -l
	1


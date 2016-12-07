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
		"owner": "7b216bd9-6523-4ca9-aa3b-1d8a5994f054",
		"identifier": "101001000281E600",
		"type": "DigitalGlobeAcquisition",
		"properties": {
			"sunElevation": "26.2667",
			"targetAzimuth": "91.47679",
			"sunAzimuth": "162.5099",
			"timestamp": "2003-11-27T00:00:00.000Z",
			"offNadirAngle": "9.0",
			"cloudCover": "6.0",
			"multiResolution": "2.480306149",
			"browseURL": "https://browse.digitalglobe.com/imagefinder/showBrowseMetadata?catalogId=101001000281E600",
			"panResolution": "0.620390713",
			"footprintWkt": "POLYGON ((-112.249662 41.26810813, -112.0471457 41.26691053, -112.0471218 41.21128254, -112.0470635 41.15576778, -112.0470193 41.1002849, -112.0468856 41.04491751, -112.0468263 40.98960299, -112.0468154 40.93429073, -112.0468277 40.87893259, -112.0467612 40.82357493, -112.0466661 40.76815214, -112.0465818 40.71267682, -112.0469022 40.65708733, -112.0468488 40.60141699, -112.0468019 40.54567617, -112.0468795 40.48981261, -112.0471183 40.43378609, -112.0466806 40.37781586, -112.0466431 40.36907251, -112.252435 40.36636078, -112.2522955 40.3751995, -112.2523337 40.4314792, -112.2516845 40.48786638, -112.2504122 40.54429583, -112.2501899 40.60029535, -112.2498266 40.65624682, -112.2495379 40.71204009, -112.2494044 40.76774333, -112.2493635 40.82337556, -112.2493227 40.87893528, -112.2492872 40.93449091, -112.2493008 40.99000069, -112.2492497 41.04548369, -112.2492948 41.10100015, -112.249356 41.15663737, -112.2495337 41.21232543, -112.249662 41.26810813))",
			"catalogID": "101001000281E600",
			"imageBands": "Pan_MS1",
			"sensorPlatformName": "QUICKBIRD02",
			"vendorName": "DigitalGlobe"
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
You can add filters for any properties in the catalog items you are searching for.  For example, here's how you return images with cloudCover
less than 50.:

.. code-block:: bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\) cloudCover \< 50


Here's a more complicated set of filters that can be applied:

.. code-block:: bash

	$ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\) offNadirAngle < 25, cloudCover < 30

Search by Types
-----------------------
You can search by type as well.  The usual type for Digital Globe Imagery is "DigitalGlobeAcquisition".  
To search only Landsat imagery for example:

.. code-block:: bash

        $ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z LandsatAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\)

Get Metadata Info about a given Catalog ID
------------------------------------------
If you have multiple catalog IDs and simply want to get records out of the catalog:

.. code-block:: bash

         $ gbdxcli catalog get 103001005DB90000 103001005BA65A00 103001005BA6CC00

Search File for a given search
------------------------------

You can redirect and reuse searches in a test file with line seperated values.

.. code-block:: bash

         $ cat aoi_search.txt

         2016-09-01T00:00:00.000Z
         2016-09-07T00:00:00.000Z
         DigitalGlobeAcquisition
         POLYGON((-113.88427734375 40.36642741921034,-110.28076171875 40.36642741921034,-110.28076171875 37.565262680889965,-113.88427734375 37.565262680889965,-113.88427734375 40.36642741921034))


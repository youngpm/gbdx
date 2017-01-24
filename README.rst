=====================================
gbdxcli: Unix binary for using GBDX
=====================================

gbdxcli is a unix cli for ordering imagery and searching the catalog on DigitalGlobe's GBDX platform.

In order to use gbdxcli, you need GBDX credentials. Email GBDX-Support@digitalglobe.com to get these.

Documentation is available in the source package on github.com.

See the license file for license rights and limitations (MIT).

Installation
------------

Copy the binary to your destination directory or build the source.  Keep in mind that the binary is compiled for Mac OSX.

If using the source, add a GOPATH environment variable to the source and build.

Get help using one of the gbdxcli --help commands.

.. code-block:: bash

         $ gbdxcli --help
         A CLI for GBDX.

         Usage:
           gbdxcli [command]

         Available Commands:
           browse      Download a browse
           catalog     Search the catalog for records using GBDX
           configure   Stores your GBDX configuration
           order       Submit orders to GBDX
           s3          Interface to S3 Storage Service
           token       Get a GBDX token
           workflow    Search for workflows using GBDX  

         Flags:
               --profile string   GBDX profile to use (default "default")

         Use "gbdxcli [command] --help" for more information about a command.

Add your GBDX credentials using gbdxcli configure.

    gbdxcli configure

Usage
---------

Here are examples of a catalog search command using command line, file redirect and line by line entering data.

.. code-block:: bash

         $ gbdxcli catalog get 103001005DB90000 103001005BA65A00 103001005BA6CC00

.. code-block:: bash

        $ gbdxcli catalog 2016-09-01T00:00:00.000Z 2016-09-07T00:00:00.000Z DigitalGlobeAcquisition POLYGON\(\(-82.7 28.945,-82.55 28.945,-82.55 28.864,-82.7 28.864,-82.7 28.945\)\) cloudCover \< 25

.. code-block:: bash

        $ gbdxcli catalog
        2017-01-22T00:00:00.000Z
        2017-01-22T23:59:59.000Z
        IDAHOImage


.. code-block:: bash

        $ gbdxcli catalog < aoi_search.txt

Development
-----------

Clone the repo::

    git clone https://github.com/youngpm/gbdx


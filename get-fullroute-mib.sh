#!/bin/bash

rm -f ./rib.*
DATE_DIRNAME=`date "+%Y.%m"`
DATE_MIBNAME=`date "+%Y%m%d.0000"`
MIBNAME=rib.${DATE_MIBNAME}.bz2

wget http://archive.routeviews.org/route-views.wide/bgpdata/${DATE_DIRNAME}/RIBS/${MIBNAME}
bzip2 -d ${MIBNAME}

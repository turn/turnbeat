#!/bin/sh

# run this to build the RPM

RPM_HOME=$HOME/rpmbuild

cp turnbeat.init $RPM_HOME/SOURCES
cp ../turnbeat $RPM_HOME/SOURCES
cp ../turnbeat.yml $RPM_HOME/SOURCES

rpmbuild -bb turnbeat.spec

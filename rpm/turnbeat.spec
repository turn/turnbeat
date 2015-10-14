Summary:	libbeat based agent with multiple functions
Name:		turnbeat
Version:        1.0
Release:        4
License:	Apache
Group:          System/Monitoring
# Distribution:   buildhash=@GIT_FULLSHA1@
Provides:	turnbeat
Source:		https://github.com/blacklightops/turnbeat

# initial super hacky rpm spec file
# expects to find following files in $RPM_SOURCE_DIR
#   turnbeat
#   turnbeat.init
#   turnbeat.yml
# you must place those there manually for now

%define __spec_prep_post true
%define __spec_prep_pre true
%define __spec_build_post true
%define __spec_build_pre true

%install
mkdir -p $RPM_BUILD_ROOT/usr/local/turnbeat
cp $RPM_SOURCE_DIR/turnbeat $RPM_BUILD_ROOT/usr/local/turnbeat

mkdir -p $RPM_BUILD_ROOT/etc/init.d
cp $RPM_SOURCE_DIR/turnbeat.init $RPM_BUILD_ROOT/etc/init.d/turnbeat
cp $RPM_SOURCE_DIR/turnbeat.yml $RPM_BUILD_ROOT/etc

%files
%dir "/usr/local/turnbeat"
"/usr/local/turnbeat/turnbeat"
"/etc/init.d/turnbeat"
"/etc/turnbeat.yml"
%attr(755, -, -) "/etc/init.d/turnbeat"

%description
libbeat based agent with multiple functions

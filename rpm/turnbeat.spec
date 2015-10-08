Summary:	libbeat based agent with multiple functions
Name:		turnbeat
Version:        1.0
Release:        1
License:	Apache
Group:          System/Monitoring
# Distribution:   buildhash=@GIT_FULLSHA1@
Provides:	turnbeat
Source:		https://github.com/turn/turnbeat

# initial super hacky rpm spec file
# expects to find turnbeat binary, and turnbeat.init in $RPM_SOURCE_DIR
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

%files
%dir "/usr/local/turnbeat"
"/usr/local/turnbeat/turnbeat"
"/etc/init.d/turnbeat"
%attr(755, -, -) "/etc/init.d/turnbeat"

%description
libbeat based agent with multiple functions

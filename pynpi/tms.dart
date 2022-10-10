// some code to do TMS tiling, should probably work ok in duckdb?

var zoom = 14;
var n = 1 << zoom;
print(n);
double lon = double.parse(test['lon']);
double lat = double.parse(test['lat']);
var xcalc = n * ((lon+180)/360);
var xtile = xcalc.floor();
var lat_rad = lat*pi/180;
print(lat_rad);
var sec_lat = 1/(cos(lat_rad));
print('sec:${sec_lat}');
var tan_lat = tan(lat_rad);
print('tan:${tan_lat}');
var y = (log(tan_lat + sec_lat))/(log(2));
print(y);
var ytrans = (1 - (y/pi))/2;
print(ytrans);
var ycalc = n * ytrans;
print(ycalc);
///var ycalc = (n * (1-((log((tan(lat_rad))+(1/(cos(lat_rad))))/log(2))/pi)))/2;
var ytile = ycalc.floor();

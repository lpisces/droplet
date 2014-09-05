awk 'BEGIN { OFS = "\t";FS = "\t"}{print $1,$2}' sz_sh_stock.TXT > stock.dat

load data local infile "/root/go-projects/src/droplet/stock.dat" into table stocks fields terminated by "\t" (code, name);

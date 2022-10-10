
#not used, but here as notes for duckdb in R
library("DBI")
con = dbConnect(duckdb::duckdb(), ":memory:")
dbWriteTable(con, "iris", iris)
dbGetQuery(con, 'SELECT "Species", MIN("Sepal.Width") FROM iris GROUP BY "Species"')

dbExecute(con, "create table npi as select * from read_csv_auto('/Users/jim/dev/asset/npi/NPPES_Data_Dissemination_September_2022/npidata_pfile_20050523-20220911.csv', header=TRUE, ALL_VARCHAR=TRUE, )")

dbExecute(con, "export database '/Users/jim/dev/asset/npi/duckdb'")
dbExecute(con, "export database '/Users/jim/dev/asset/npi/duckdb2' (FORMAT PARQUET) ")

dbGetQuery(con, 'select distinct from npi where "Provider First Name"<>')
dbExecute(con, "drop table npi")

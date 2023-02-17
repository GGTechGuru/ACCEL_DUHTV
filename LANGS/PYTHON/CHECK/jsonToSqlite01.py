# Json from file
# Get table names (outside key) and schema (value)
# Get basic name/types & add table gen statements
# Get compound name-><X>
# If compound value -> {<type>: {CHOICES: [<list>] }, add name/type/check(name=... or...)
# If compound value -> {<type>: {SOURCE_DATA:{...}}:
#   Add (a) name/type
#   Add (b) foreign key <references>

import json
import sqlite3
import sys

class TYPES:
    STR = "<class 'str'>"
    DICT = "<class 'dict'>"
    LIST = "<class 'list'>"

    TYPES_LIST = [ STR, DICT, LIST ]

    @staticmethod
    def istype(val):
        type_str = str(type(val))

        print( type_str )

        if (type_str in TYPES.TYPES_LIST):
            return type_str
        else:
            return None

###################

class JsonToSqlite:

    def get_table_schemas(self, schemas_json_file):
        json_str = (open(schemas_json_file)).read()
        schemas_json = json.loads(json_str)

        table_schemas = schemas_json["TABLE_SCHEMAS"]

        # for name in table_schemas.keys():
            # print(name)

        return table_schemas

    ##############################################

    def schema_to_ddl(self, tbl_name, table_schema):
        print(tbl_name)
        print(table_schema.values())

        # SQLite Syntax?
        crt_tbl_ddl = "CREATE TABLE IF NOT EXISTS %s ( " % (tbl_name)

        count = 0

        for column in table_schema.keys():
            col_descr = table_schema[column]
            print(str(type(col_descr)))

            if count > 0:
                crt_tbl_ddl += ", "

            count += 1

            if TYPES.istype(col_descr) == TYPES.STR:
                col_ddl = " %s %s " % (column, col_descr)
                print(col_ddl)
                crt_tbl_ddl += col_ddl

            # -> compound
            elif TYPES.istype(col_descr) == TYPES.DICT:
                type_descr = (list(col_descr.keys()))[0]

                inner_descr = col_descr[type_descr]

                # constraint_descr = ((list(inner_descr.keys()))[0]).upper()
                constraint_descr = (list(inner_descr.keys()))[0]

                # -> foreign key / references other table+col
                if constraint_descr.upper() == "SOURCE_DATA":

                    tbl_col_descr = inner_descr[constraint_descr]

                    ref_tbl_name = (list(tbl_col_descr.keys()))[0]
                    ref_col_name = tbl_col_descr[ref_tbl_name]

                    col_ddl = " %s %s , " % (column, type_descr)
                    print(col_ddl)
                    crt_tbl_ddl += col_ddl

                    col_ddl = " FOREIGN KEY (%s) REFERENCES %s(%s) " % (column, ref_tbl_name, ref_col_name)
                    print(col_ddl)

                    crt_tbl_ddl += col_ddl


                # -> choices/list check
                elif constraint_descr.upper() == "CHOICES":
                    choice_vals = inner_descr[constraint_descr]
                    if TYPES.istype(choice_vals) == TYPES.LIST:

                        opt_quotes = ''
                        if type_descr.upper() in [ "STRING", "VARCHAR", "TEXT" ]: # Varchar?
                            opt_quotes = '"'

                        col_ddl = " %s %s " % ( column, type_descr )
                        col_ddl += " CHECK ( "

                        ok_count = 0
                        for ok_value in choice_vals:

                            if ok_count > 0:
                                col_ddl += " OR "

                            col_ddl += (" %s = %s%s%s " % (column, opt_quotes, str(ok_value), opt_quotes))

                            ok_count += 1

                        col_ddl += " ) "

                        crt_tbl_ddl += col_ddl

                    else:
                        err_msg = "Choice values should be a list [%s]" % (choice_vals)
                        raise Exception(err_msg)

                # -> Need to handle
                else:
                    err_msg = "Unknown constraint type [%s]" % (constraint_descr)
                    raise Exception(err_msg)


        crt_tbl_ddl += " ); " # Syntax?

        print(crt_tbl_ddl)

        return crt_tbl_ddl
    
#######################################################################

# main

schemas_json_file = None
if len(sys.argv) < 2:
    print("Need 1 argument: JSON schema file path")
    sys.exit(1)
else:
    schemas_json_file = sys.argv[1]

jts = JsonToSqlite()

table_schemas = jts.get_table_schemas(schemas_json_file)

tbl_name = (list(table_schemas.keys()))[0]
test_json = table_schemas[tbl_name]

jts.schema_to_ddl(tbl_name, test_json)



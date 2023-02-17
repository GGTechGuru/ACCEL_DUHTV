package main

// Convert JSON definition of SQLite tables into default read (all columns, no conditions)
// and default insert (all columns, exclude auto-incrementing keys [defer]).

// CLI args:
//	(1) Output Golang file path
//	(2) JSON SQLite DB definition: file path

// Helper functions:
//	+ JSON type -> Golang type
//	+ JSON/Golang type -> printf % code


// Do:
//    Parse column and type names from JSON
//    Build standard function string with
//        Corresponding Golang types
//        Corresponding printf % codes
//
//    Write to Golang SQLite access file

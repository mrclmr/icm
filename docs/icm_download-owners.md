## icm download-owners

Download information of owners and write CSV to file

### Synopsis

Download information of owners and write CSV to file.
Following information is available:

  Owner code
  Company
  City
  Country

```
icm download-owners [flags]
```

### Examples

```
# Overwrite owner.csv file with newest owners
icm download-owners
# Create custom-owner.csv to have additional custom mapping of owner codes
# Use semicolon as a separator. For using double quotes please see existing
# owner.csv file.
echo 'AAA;my company;my city;my country' >> $HOME/.icm/data/custom-owner.csv
```

### Options

```
  -h, --help            help for download-owners
  -o, --output string   output file (default "$HOME/.icm/data/owner.csv")
```

### SEE ALSO

* [icm](icm.md)	 - Validate or generate intermodal container markings


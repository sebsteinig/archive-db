# archive-db
## API Routes
- Search : /search/?like={prefix}&{...}  
For example :
```console
http://localhost:3000/search/?like=te&extension=png
```
Additional parameters can be : config_name, extension, lossless, threshold, rx, ry and chunks. They can be specified in any order. 
- Select :  /select/{exp_id}/?{parameters=...}  
For example :
```console
http://localhost:3000/select/texpa1/?extension=png&variables=[tos,winds]&lossless=true
```
Additional parameters can be : variables, config_name, extension, lossless, threshold, rx, ry and chunks. They can be specified in any order.
If the parameter variables is specified, it has to be a list written with brackets or there will be an error.
- Select collection : /select/collection/?ids=[{exp_ids}]&{parameters=...}  
For example :
```console
http://localhost:3000/select/collection/?ids=[texpa1,xooeh]&extension=png&variables=[tos,winds]&lossless=true
```
This request is like a normal select. If the ids are not specified, there will be an error. Just like the variables parameter, the list of ids has to be written under brackets.
- Insert : route used by Nimbus.

import { IsArray, IsBoolean, IsNotEmpty, IsObject, IsOptional, IsString } from "class-validator"
import { ConversionInfo } from "./conversion_info.dto"

export class Experiment {
    @IsString()
    @IsNotEmpty()
    exp_id : string

    @IsString()
    @IsNotEmpty()
    config : string

    @IsString()
    @IsNotEmpty()
    nimbus_version : string

    @IsBoolean()
    var_clt               : boolean = false
    @IsBoolean()
    var_currents          : boolean = false
    @IsBoolean()
    var_height            : boolean = false
    @IsBoolean()
    var_liconc            : boolean = false
    @IsBoolean()
    var_mlosts            : boolean = false
    @IsBoolean()
    var_pfts              : boolean = false
    @IsBoolean()
    var_pr                : boolean = false
    @IsBoolean()
    var_sic               : boolean = false
    @IsBoolean()
    var_snc               : boolean = false
    @IsBoolean()
    var_tas               : boolean = false
    @IsBoolean()
    var_tos               : boolean = false
    @IsBoolean()
    var_winds             : boolean = false
    
    @IsOptional()
    @IsObject()
    metadata              : Object 

    @IsArray()
    @IsNotEmpty()
    conversion_info       : ConversionInfo[]
}

export interface ImagePath {
    path : String,
    extension : String,
    metadata : Object
}

export interface ConversionInfo {
    levels : number,
    time_steps : number,
    xsize : number,
    ysize : number,              
    xfirst : number,
    xinc : number,
    yfirst : number,
    yinc : number
    nan_value_encoding : number,
    threshold : number,
    image_paths : ImagePath[],
}

export interface Experiment {
    exp_id : String,
    config : String,
    nimbus_version : String,
    var_clt               : boolean,
    var_currents          : boolean,
    var_height            : boolean,
    var_liconc            : boolean,
    var_mlosts            : boolean,
    var_pfts              : boolean,
    var_pr                : boolean,
    var_sic               : boolean,
    var_snc               : boolean,
    var_tas               : boolean,
    var_tos               : boolean,
    var_winds             : boolean,
    metadata              : Object ,
    conversion_info       : ConversionInfo[],
}
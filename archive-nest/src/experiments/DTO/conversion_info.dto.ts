import { IsArray, IsInt, IsNotEmpty, IsNumber, IsPositive } from "class-validator"
import { ImagePath } from "./image_path.dto"


export class ConversionInfo {
    @IsNumber()
    @IsPositive()
    @IsInt()
    @IsNotEmpty()
    levels : number


    @IsNumber()
    @IsPositive()
    @IsInt()
    @IsNotEmpty()
    time_steps : number


    @IsNumber()
    @IsPositive()
    @IsInt()
    @IsNotEmpty()
    xsize : number


    @IsNumber()
    @IsPositive()
    @IsInt()
    @IsNotEmpty()
    ysize : number  
    
    
    @IsNumber()
    @IsNotEmpty()
    xfirst : number


    @IsNumber()
    @IsNotEmpty()
    xinc : number


    @IsNumber()
    @IsNotEmpty()
    yfirst : number


    @IsNumber()
    @IsNotEmpty()
    yinc : number


    @IsNumber()
    @IsPositive()
    @IsInt()
    @IsNotEmpty()
    nan_value_encoding : number


    @IsNumber()
    @IsNotEmpty()
    threshold : number

    @IsArray()
    @IsNotEmpty()
    image_paths : ImagePath[]
}
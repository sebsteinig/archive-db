import { IsNotEmpty, IsObject, IsOptional, IsString } from "class-validator"

export class ImagePath {
    @IsString()
    @IsNotEmpty()
    path : string

    @IsString()
    @IsNotEmpty()
    extension : string

    @IsObject()
    @IsOptional()
    metadata : Object
}
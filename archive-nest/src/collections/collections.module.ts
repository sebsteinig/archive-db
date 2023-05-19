import { Global, Module } from "@nestjs/common";
import { CollectionsService } from "./collections.service";
import { CollectionsController } from "./collections.controller";
import { PrismaModule } from "src/prisma/prisma.module";

@Global()
@Module({
    controllers : [CollectionsController],
    providers : [CollectionsService],
    exports : [CollectionsService]
})
export class CollectionsModule{}
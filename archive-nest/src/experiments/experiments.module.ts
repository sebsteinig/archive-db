import { Global, Module } from "@nestjs/common";
import { ExperimentsService } from "./experiments.service";
import { ExperimentsController } from "./experiments.controller";
import { PrismaModule } from "src/prisma/prisma.module";

@Global()
@Module({
    controllers : [ExperimentsController],
    providers : [ExperimentsService],
    exports : [ExperimentsService]
})
export class ExperimentsModule{}
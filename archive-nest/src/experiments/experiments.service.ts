import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { Experiment } from './DTO';

@Injectable()
export class ExperimentsService {

    constructor(private prisma : PrismaService) {}

    insertExperiment(experiment : Experiment) {
        //this.prisma.
    }

}

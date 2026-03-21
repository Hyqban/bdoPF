import { Component, Input } from '@angular/core';
import { SearchService } from '../../services/search-service';

@Component({
    selector: 'app-amount',
    imports: [],
    templateUrl: './amount.html',
    styleUrl: './amount.scss',
})
export class Amount {
    constructor(protected search: SearchService) {}

    @Input() breadCrumbLength: number = 0;
    @Input() itemCount: number = 1;
}

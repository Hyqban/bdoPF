import {
    ChangeDetectionStrategy,
    Component,
    inject,
    Input,
    OnInit,
    signal,
    Signal,
    WritableSignal,
} from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { Item, RecipeAmountInterface } from '../../shared/models/model';
import { SearchService } from '../../services/search-service';
import { debounceTime } from 'rxjs';

@Component({
    selector: 'app-recipe-amount',
    imports: [MatIconModule, MatIconModule, ReactiveFormsModule],
    templateUrl: './recipe-amount.html',
    styleUrl: './recipe-amount.scss',
})
export class RecipeAmount implements OnInit {
    // @Input() recipeAmount!: WritableSignal<RecipeAmountInterface>;
    // items = signal<cloneDeep(this.recipeAmount().items);
    search = inject(SearchService);
    // items = signal(this.search.recipeAmount().items);
    form!: FormGroup;

    ngOnInit(): void {
        this.form = new FormGroup({
            requiredInput: new FormControl(''),
            optionalInput: new FormControl(''),
        });

        this.form.valueChanges.pipe(debounceTime(500)).subscribe((value) => {
            if (value.requiredInput === '') {
                this.search.recipeAmount.update((el) => ({
                    ...el,
                    amount: 1,
                }));
            }
            if (value.optionalInput === '') {
                this.search.recipeAmount.update((el) => ({
                    ...el,
                    averageYield: 1,
                }));
            }

            if (Number(value.requiredInput) >= 1) {
                this.search.recipeAmount.update((el) => ({
                    ...el,
                    amount: Math.floor(Number(value.requiredInput)),
                }));
            }
            if (Number(value.optionalInput) >= 0.1) {
                this.search.recipeAmount.update((el) => ({
                    ...el,
                    averageYield: Number(value.optionalInput),
                }));
            }

            this.search.recipeAmount.update((el) => ({
                ...el,
                amountItems: el.items,
            }));

            this.search.recipeAmount.update((ra) => ({
                ...ra,
                amountItems: ra.amountItems.map((el) => ({
                    ...el,
                    count: el.count
                        ? Math.ceil((el.count * ra.amount) / ra.averageYield)
                        : el.count,
                })),
            }));
        });
    }

    close() {
        this.search.recipeAmount.update((el) => ({
            ...el,
            open: false,
            items: [],
            amountItems: [],
            amount: 1,
            averageYield: 1,
        }));
    }
}

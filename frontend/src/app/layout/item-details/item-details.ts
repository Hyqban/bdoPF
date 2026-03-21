import { ChangeDetectionStrategy, Component, signal, WritableSignal } from '@angular/core';
import { SearchService } from '../../services/search-service';
import { MatCardModule } from '@angular/material/card';
import { MatTabsModule } from '@angular/material/tabs';
import { I18nService } from '../../services/i18n-service';
import { Amount } from '../amount/amount';
import { MatIconModule } from '@angular/material/icon';
import { RecipeAmount } from '../recipe-amount/recipe-amount';
import { Item, RecipeAmountInterface } from '../../shared/models/model';

@Component({
    selector: 'app-item-details',
    imports: [MatCardModule, MatTabsModule, MatIconModule, Amount, RecipeAmount],
    templateUrl: './item-details.html',
    styleUrl: './item-details.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ItemDetails {
    constructor(
        protected search: SearchService,
        protected i18n: I18nService,
    ) {}

    // recipeAmount: WritableSignal<RecipeAmountInterface> = signal<RecipeAmountInterface>({
    //     open: false,
    //     items: [],
    //     amount: 1,
    //     averageYield: 1,
    // });

    openRecipeAmount(items: Item[]) {
        this.search.recipeAmount.update((el) => ({
            ...el,
            open: true,
            items: items,
            amountItems: items,
        }));
        // this.recipeAmount.update((el) => ({
        //     ...el,
        //     open: true,
        //     items: items,
        // }));
    }
}

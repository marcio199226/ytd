import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { FormControl } from '@angular/forms';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class HomeComponent implements OnInit {
  public searchInput: FormControl;

  constructor() {
    this.searchInput = new FormControl('');
  }

  ngOnInit(): void {}

  clearSearch(): void {
    this.searchInput.setValue('');
  }
}

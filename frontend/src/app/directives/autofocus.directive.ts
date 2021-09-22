import { Directive, OnInit, ElementRef } from '@angular/core';

@Directive({
  selector: '[autofocus]'
})
export class AutofocusDirective implements OnInit {

  constructor(private elementRef: ElementRef) { };

  ngOnInit(): void {
    this.elementRef.nativeElement.focus();
  }

}

import { this$ } from 'sham-ui-macro/ref.macro';

export function SetField( options, update ) {
    this$.setField = field => ( value ) => update( {
        [ field ]: value
    } );
}

export function ToggleField( options, update ) {
    const state = options();
    this$.toggleField = field => () => update( {
        [ field ]: !state[ field ]
    } );
}

import * as directives from 'sham-ui-directives';
import FormContainer  from '../../../src/components/FormContainer.sfc';
import renderer  from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( FormContainer, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'submit', async() => {
    expect.assertions( 1 );
    const success = jest.fn();
    const meta = renderer( FormContainer, {
        directives,
        success
    } );
    meta.component.container.querySelector( '[type="submit"]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( success ).toHaveBeenCalledTimes( 1 );
} );

it( 'submit fail', async() => {
    expect.assertions( 2 );
    const submit = jest.fn().mockReturnValueOnce( Promise.reject() );
    const meta = renderer( FormContainer, {
        directives,
        submit
    } );
    meta.component.container.querySelector( '[type="submit"]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( submit ).toHaveBeenCalledTimes( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

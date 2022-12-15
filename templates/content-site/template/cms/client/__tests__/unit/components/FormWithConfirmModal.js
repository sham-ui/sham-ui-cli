import * as directives from 'sham-ui-directives';
import FormWithConfirmModal  from '../../../src/components/FormWithConfirmModal.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( FormWithConfirmModal, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'success save', async() => {
    expect.assertions( 3 );
    const onSuccess = jest.fn();
    const meta = renderer( FormWithConfirmModal, {
        onSuccess
    }, {
        directives
    } );

    meta.ctx.container.querySelector( '[type="submit"]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();

    meta.ctx.container.querySelector( '[data-test-ok-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );

    expect( onSuccess ).toHaveBeenCalledTimes( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'fail save', async() => {
    expect.assertions( 1 );
    const meta = renderer( FormWithConfirmModal, {
        saveData: () => Promise.reject( [ 'Error test' ] )
    }, {
        directives
    } );

    meta.ctx.container.querySelector( '[type="submit"]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );

    meta.ctx.container.querySelector( '[data-test-ok-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );

    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'cancel save', async() => {
    expect.assertions( 2 );
    const saveData = jest.fn();
    const meta = renderer( FormWithConfirmModal, {
        saveData
    }, {
        directives
    } );

    meta.ctx.container.querySelector( '[type="submit"]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );

    meta.ctx.container.querySelector( '[data-test-cancel-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );

    expect( saveData ).toHaveBeenCalledTimes( 0 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

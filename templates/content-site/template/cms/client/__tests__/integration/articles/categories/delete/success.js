import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success delete cateogory', async() => {
    expect.assertions( 4 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
            ]
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    await app.waitRendering();

    app.checkMainPanel();

    axios
        .use( 'delete', '/categories/1' )
        .use( 'get', '/categories', {
            categories: [
                { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
            ]
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.delete.mock.calls[ 0 ][ 0 ] ).toBe( '/categories/1' );
    app.checkMainPanel();
} );

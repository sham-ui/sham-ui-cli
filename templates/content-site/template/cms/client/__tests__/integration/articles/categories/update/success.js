import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );
const DI = setup();

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success update category data', async() => {
    expect.assertions( 5 );

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

    app.click( '[data-test-update-button="6"]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        name: 'New name'
    };

    app.form.fill( 'name', formData.name );
    app.form.submit();

    axios
        .use( 'put', '/categories/6', {} )
        .use( 'get', '/categories', {
            categories: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'New name', 'slug': 'new-name' }
            ]
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( '/categories/6' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkMainPanel();
} );

import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success delete tag', async() => {
    expect.assertions( 4 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/tags', {
            tags: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    await app.waitRendering();

    app.checkMainPanel();

    axios
        .use( 'delete', '/tags/1' )
        .use( 'get', '/tags', {
            tags: [
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.delete.mock.calls[ 0 ][ 0 ] ).toBe( '/tags/1' );
    app.checkMainPanel();
} );

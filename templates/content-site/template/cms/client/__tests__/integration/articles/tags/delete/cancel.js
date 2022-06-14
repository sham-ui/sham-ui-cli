import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel delete tag', async() => {
    expect.assertions( 2 );

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
    app.click( '[data-test-modal] [data-test-cancel-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 0 );
    app.checkMainPanel();
} );

import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success create category', async() => {
    expect.assertions( 5 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: []
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        name: 'Travel'
    };

    axios
        .use( 'post', '/categories', {} )
        .use( 'get', '/categories', {
            categories: [
                {
                    id: 1,
                    name: formData.name,
                    slug: 'travel'
                }
            ]
        } );

    app.form.fill( 'name', formData.name );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/categories' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
} );

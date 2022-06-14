import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import RoutesLoginPage  from '../../../../../src/components/routes/login/page.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';


it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );

    const meta = renderer( RoutesLoginPage, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );

    const loginMock = jest.fn();
    DI.bind( 'session', {
        login: loginMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesLoginPage, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );

    const formData = {
        email: 'admin@gmail.com',
        password: 'passw0rd'
    };
    const { component } = meta;
    component.container.querySelector( '[name="email"]' ).value = formData.email;
    component.container.querySelector( '[name="password"]' ).value = formData.password;
    component.container.querySelector( '[type="submit"]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( loginMock.mock.calls ).toHaveLength( 1 );
    expect( loginMock.mock.calls[ 0 ] ).toHaveLength( 2 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

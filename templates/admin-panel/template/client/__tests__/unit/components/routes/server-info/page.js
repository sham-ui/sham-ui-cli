import RoutesServerInfoPage  from '../../../../../src/components/routes/server-info/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import { storage } from '../../../../../src/storages/session';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        api: {
            request: getMock.mockReturnValueOnce(
                Promise.resolve( { } )
            )
        }
    } );

    const meta = renderer( RoutesServerInfoPage, {
        DI,
        directives,
        filters: {}
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'store', {
        api: {
            request: jest.fn().mockReturnValueOnce(
                Promise.reject( {} )
            )
        }
    } );
    const meta = renderer( RoutesServerInfoPage, {
        DI,
        directives,
        filters: {}
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

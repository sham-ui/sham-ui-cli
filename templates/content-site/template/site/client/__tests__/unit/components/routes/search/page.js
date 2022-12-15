import * as directives from 'sham-ui-directives';
import { createDI } from 'sham-ui';
import RoutesSearchPage  from '../../../../../src/components/routes/search/page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );

    DI.bind( 'store', {
        articles: jest.fn().mockReturnValueOnce(
            Promise.resolve( { meta: {} } )
        )
    } );

    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            }
        },
        generate: () => '/'
    } );
    DI.bind( 'location:href', 'http://client.example.com/search/1/?q=content' );

    const meta = renderer( RoutesSearchPage, {}, {
        DI,
        directives
    } );

    await new Promise( resolve => setImmediate( resolve ) );

    expect( meta.toJSON() ).toMatchSnapshot();
} );

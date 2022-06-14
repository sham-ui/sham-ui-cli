import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import renderer from 'sham-ui-test-helpers';
//eslint-disable-next-line max-len
import RoutesHomePagination  from '../../../../../../src/components/routes/_common/articles/pagination.sfc';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            },
            url: '',
            addWatcher() {}
        },
        generate: () => '/'
    } );

    const meta = renderer( RoutesHomePagination, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

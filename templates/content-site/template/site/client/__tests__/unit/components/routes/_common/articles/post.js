import { createDI } from 'sham-ui';
import hrefto from 'sham-ui-router/lib/href-to';
import setupUnsafe from 'sham-ui-unsafe';
import formatLocaleDate from '../../../../../../src/filters/format-locale-date';
import RoutesHomePost  from '../../../../../../src/components/routes/_common/articles/post.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'router', {
        generate: () => '/'
    } );

    setupUnsafe( DI );

    const meta = renderer( RoutesHomePost, {
        DI,
        directives: {
            hrefto
        },
        filters: {
            formatLocaleDate
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

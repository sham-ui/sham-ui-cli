// eslint-disable-next-line max-len
import RoutesArticlesCommonFieldsSelect  from '../../../../../../../src/components/routes/articles/common/fields/select.sfc';
import { oninput } from 'sham-ui-directives';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCommonFieldsSelect, {
        directives: {
            oninput
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

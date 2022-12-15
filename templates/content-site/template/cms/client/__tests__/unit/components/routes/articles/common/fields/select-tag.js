import * as directives from 'sham-ui-directives';
import { onFocusIn } from '../../../../../../../src/directives/on-focus-in';
// eslint-disable-next-line max-len
import RoutesArticlesCommonFieldsSelectTag  from '../../../../../../../src/components/routes/articles/common/fields/select-tag.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCommonFieldsSelectTag, {}, {
        directives: {
            onFocusIn,
            ...directives
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

/**
 * Mini Elastic Sift. Email client controller entry point.
 */
import { EmailClientController, registerEmailClientController } from '@redsift/sift-sdk-web';

export default class MyEmailClientController extends EmailClientController {
  constructor() {
    super();
  }

  // for more info: http://docs.redsift.com/docs/client-code-redsiftclient
  loadThreadListView (listInfo) {
    console.log('mini-elastic: loadThreadListView: ', listInfo);
    // if (listInfo) {
    //   return {
    //     template: '001_list_common_txt',
    //     value: {
    //       color: '#ffffff',
    //       backgroundColor: '#e11010',
    //       subtitle: 'subtitle'
    //     }
    //   };
    // }
  };
}

// Do not remove. The Sift is responsible for registering its views and controllers
registerEmailClientController(new MyEmailClientController());

/**
 * Mini Elastic Sift. Frontend view entry point.
 */
import { SiftView, registerSiftView } from '@redsift/sift-sdk-web';

const nanoToMillis = v => Math.round(v/Math.pow(10, 4)) / 100

export default class MyView extends SiftView {
  constructor() {
    // You have to call the super() method to initialize the base class.
    super();
    this.onStatsUpdate = this.onStatsUpdate.bind(this);
  }

  // for more info: http://docs.redsift.com/docs/client-code-siftview
  presentView(value) {
    console.log('mini-elastic: presentView: ', value);
    document.querySelector("#api_url").innerHTML = value.data.apiUrl;
    this.controller.subscribe('stats', this.onStatsUpdate);
    this.onStatsUpdate(value.data.stats)
  };

  willPresentView(value) { };


  onStatsUpdate(stats) {
    if (!stats || !stats.index) {
      return;
    }
    ["analysis_time", "index_time"].forEach(x => stats.index[x] = nanoToMillis(stats.index[x]) +'ms' );
    stats["search_time"] = nanoToMillis(stats["search_time"]) +'ms'
    document.querySelector("#index_stats").innerHTML = JSON.stringify(stats, null, ' ')
  }

}

registerSiftView(new MyView(window));

use warp::http::StatusCode;
use warp::reply;
use warp::Filter;

#[tokio::main]
async fn main() {
    let two_sum = warp::path!(u32 / u32).map(|a, b| format!("{}", a + b));

    let two_path = warp::path!(String / String).map(|a, b| {
        let reply = format!("Hello from warp/{}/{}! ", a, b);
        reply::with_status(reply, StatusCode::CREATED)
    });

    let root = two_sum.or(two_path);

    println!("[log] Running warp application!");
    warp::serve(root).run(([127, 0, 0, 1], 3000)).await;
}

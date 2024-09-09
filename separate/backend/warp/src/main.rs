use warp::Filter;

#[tokio::main]
async fn main() {
    let root = warp::any()
        .and(warp::path::param::<String>())
        .map(|s| format!("Hello from warp/{}! ", s))
        .and(warp::path::param::<String>())
        .map(|a, b| format!("Hello from warp/{}/{}! ", a, b));

    println!("[log] Running warp application!");
    warp::serve(root).run(([127, 0, 0, 1], 3000)).await;
}

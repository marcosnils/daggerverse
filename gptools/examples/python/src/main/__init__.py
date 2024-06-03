import dagger
from dagger import dag, function, object_type


@object_type
class Python:
    @function
    async def gptools_transcript__directory(
        self,
        url: str = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
    ) -> dagger.Directory:
        """example on how to get a transcript for a video and
        return it as a *Directory to use in other functions"""
        # TODO: call transcript once python gen is fixed
        # video = dag.http(url)
        return dag.directory().with_file(
            "video-transcript.txt", dag.gptools().audio(url)
        )
